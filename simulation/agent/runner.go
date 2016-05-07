package agent

// Runner must be implemented for every agent
type Runner interface {
	Run() *Chans
}
