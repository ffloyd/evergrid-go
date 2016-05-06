package agent

// Runner must be implemented for interaction with GlobalTimer
type Runner interface {
	Run() *Chans
}
