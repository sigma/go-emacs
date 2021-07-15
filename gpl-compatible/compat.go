// Package gpl_compatible is a marker package, meant to be imported for side-effect only.
// It will export the expected "plugin_is_GPL_compatible" symbol in an Emacs module,
// and acknowledges that the importing package indeed is GPL compatible.
package gpl_compatible

// int plugin_is_GPL_compatible;
import "C"
