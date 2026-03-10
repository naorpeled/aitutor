package types

// Tier represents the difficulty level of a lesson.
type Tier int

const (
	Beginner     Tier = iota
	Intermediate
	Advanced
)

func (t Tier) String() string {
	switch t {
	case Beginner:
		return "Beginner"
	case Intermediate:
		return "Intermediate"
	case Advanced:
		return "Advanced"
	default:
		return "Unknown"
	}
}

// BlockKind represents the type of a theory block.
type BlockKind int

const (
	Paragraph BlockKind = iota
	Heading
	Code
	Callout
	Bullet
)

// TheoryBlock is a single block of content in the theory phase.
type TheoryBlock struct {
	Kind    BlockKind
	Content string
}

// QuizKind represents the type of a quiz question.
type QuizKind int

const (
	MultipleChoice QuizKind = iota
	FillBlank
	Ordering
)

// QuizQuestion defines a single quiz question.
type QuizQuestion struct {
	Kind        QuizKind
	Prompt      string
	Choices     []string
	CorrectIdx  int
	Answer      string
	Explanation string
}

// LessonDef defines a complete lesson.
type LessonDef struct {
	ID         int
	Title      string
	Tier       Tier
	Summary    string
	Theory     []TheoryBlock
	VizBuilder func(w, h int) interface{} // returns a viz.Model
	Questions  []QuizQuestion
}
