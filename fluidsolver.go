package manta

import "math"

type FluidSolver struct {
	Dt           float64
	LockDt       bool
	TimePerFrame float64
	TimeTotal    float64
	Frame        int
	FrameLength  float64
	CflCond      float64
	DtMin        float64
	DtMax        float64
	GridSize     vec3[int64]
	Dim          int
	FourthDim    bool
}

func NewFluidSolver(gridSize vec3[int64], dim int, fourthDim bool) *FluidSolver {
	if dim == 4 && !fourthDim {
		panic("Don't create 4D solvers, use 3D with fourthDim parameter = true instead.")
	}
	if dim != 2 && dim != 3 {
		panic("Only 2D and 3D solvers allowed.")
	}
	return &FluidSolver{
		Dt:           1,
		LockDt:       false,
		TimePerFrame: 0,
		TimeTotal:    0,
		Frame:        0,
		FrameLength:  1,
		CflCond:      1000,
		DtMin:        1,
		DtMax:        1,
		GridSize:     gridSize,
		Dim:          dim,
		FourthDim:    fourthDim,
	}
}

func (f *FluidSolver) Step() {
	f.TimePerFrame += f.Dt
	f.TimeTotal += f.Dt
	if f.TimePerFrame < VectorEpsilon {
		f.Frame++
		// re-calc total time, prevent drift
		f.TimeTotal = float64(f.Frame) * f.FrameLength
		f.TimePerFrame = 0
		f.LockDt = false
	}
}

// warning, uses 10^-4 epsilon values, thus only use around "regular" FPS time scales, e.g. 30
// frames per time unit pass max magnitude of current velocity as maxvel, not yet scaled by dt!
func (f *FluidSolver) AdaptTimestep(maxVel float64) {
	mvt := maxVel * f.Dt
	if !f.LockDt {
		// calculate current timestep from maxvel, clamp range
		f.Dt = math.Max(math.Min(f.Dt*f.CflCond/(mvt+1e-05), f.DtMax), f.DtMin)
		if (f.TimePerFrame + f.Dt*1.05) > f.FrameLength {
			// within 5% of full step? add epsilon to prevent roundoff errors
			f.Dt = f.FrameLength - f.TimePerFrame + 1e-04
		} else if (f.TimePerFrame+f.Dt+f.DtMin) > f.FrameLength || (f.TimePerFrame+f.Dt*1.25) > f.FrameLength {
			// avoid tiny timesteps and strongly varying ones, do 2 medium size ones if necessary
			f.Dt = (f.FrameLength - f.TimePerFrame + 1e-04) * 0.5
			f.LockDt = true
		}
	}
	if f.Dt <= f.DtMin/2 {
		panic("Invalid dt encountered! Shouldn't happen.")
	}
}
