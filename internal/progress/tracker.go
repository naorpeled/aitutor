package progress

// Tracker manages in-memory progress state.
type Tracker struct {
	data  ProgressData
	total int
	dirty bool
}

// NewTracker creates a tracker, loading saved progress.
func NewTracker(totalLessons int) *Tracker {
	data, _ := Load()
	return &Tracker{
		data:  data,
		total: totalLessons,
	}
}

// CompleteLesson marks a lesson as completed and saves.
func (t *Tracker) CompleteLesson(lessonID int) {
	t.data.CompletedLessons[lessonID] = true
	t.dirty = true
	t.save()
}

// SetLastLesson records which lesson the user was on.
func (t *Tracker) SetLastLesson(idx int) {
	t.data.LastLessonIdx = idx
	t.dirty = true
	t.save()
}

// IsCompleted returns whether a lesson is completed.
func (t *Tracker) IsCompleted(lessonID int) bool {
	return t.data.CompletedLessons[lessonID]
}

// CompletedCount returns how many lessons are completed.
func (t *Tracker) CompletedCount() int {
	return len(t.data.CompletedLessons)
}

// Total returns the total number of lessons.
func (t *Tracker) Total() int {
	return t.total
}

// LastLessonIdx returns the last viewed lesson index.
func (t *Tracker) LastLessonIdx() int {
	return t.data.LastLessonIdx
}

// CompletedMap returns the completed lessons map (for sidebar).
func (t *Tracker) CompletedMap() map[int]bool {
	return t.data.CompletedLessons
}

func (t *Tracker) save() {
	_ = Save(t.data)
	t.dirty = false
}
