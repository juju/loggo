// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

// WriterNames returns the names of the context's writers for testing purposes.
func (c *Context) WriterNames() []string {
	c.writersMutex.Lock()
	defer c.writersMutex.Unlock()
	var result []string
	for name := range c.writers {
		result = append(result, name)
	}
	return result
}

func ResetDefaultContext() {
	ResetLogging()
	DefaultContext().AddWriter(DefaultWriterName, defaultWriter())
}
