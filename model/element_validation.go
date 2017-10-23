package model

// ElementValidation
// Inspired by Joi https://github.com/hapijs/joi
type ElementValidation struct {
	// Required means the element must not be a zero value (0, false, "")
	Required bool

	// MustProvide means the element must be submitted on every POST, PUT AND PATCH
	MustProvide        bool
	MustProvideOnPost  bool
	MustProvideOnPut   bool
	MustProvideOnPatch bool

	//Forbidden bool // must not be sent
	//Strip bool  // remove from output after validation
	//Any bool  // any data type
	//String bool
	//Min int
	//Max int
	//Allow string // whitelist
	//Disallow string // blacklist
	// TODO ... loads more from Joi ...
}
