package worklist

// Entry is a struct with a single field, Path, which is a string.
// @property {string} Path - The path to the file or directory.
type Entry struct {
	Path string
}

// @property jobs - a channel of Entry
type Worklist struct {
	jobs chan Entry
}

// Adding a job to the channel.
func (w *Worklist) Add(work Entry) {
	w.jobs <- work
}

// Receiving a value from the channel and returning it.
func (w *Worklist) Next() Entry {
	j := <-w.jobs
	return j
}

// `New` creates a new worklist with a buffer size of `bufSize`
func New(bufSize int) Worklist {
	return Worklist{make(chan Entry, bufSize)}
}

// `NewJob` is a function that takes a string and returns an Entry
func NewJob(path string) Entry {
	return Entry{path}
}

// Adding a blank entry to the channel. (SELF TERMINATE)
func (w *Worklist) Finalize(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w.Add((Entry{""}))
	}
}
