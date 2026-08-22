package advanced

import (
	"github.com/naorpeled/aitutor/internal/lesson"
	"github.com/naorpeled/aitutor/internal/viz"
	"github.com/naorpeled/aitutor/pkg/types"
)

func init() {
	lesson.Register(types.LessonDef{
		ID:         18,
		Title:      "Building Reliable Agents",
		Tier:       types.Advanced,
		Summary:    "Compaction, memory storage, retrieval, and RAG",
		SourceFile: "internal/content/advanced/18_agents_building.go",
		VizBuilder: func(w, h int) interface{} { return viz.NewAgentBuildingModel(w, h) },
		Theory: []types.TheoryBlock{
			{Kind: types.Heading, Content: "Building Reliable Agents"},
			{Kind: types.Paragraph, Content: "An AI coding agent is more than a prompt wrapped around tools. Reliable agents manage state deliberately: what stays in the active context window, what gets compacted into a handoff summary, and what becomes long-term memory for future sessions."},
			{Kind: types.Heading, Content: "Compaction"},
			{Kind: types.Paragraph, Content: "Compaction is the process of summarizing a long working session into a smaller, high-signal state bundle. A good compact summary keeps the goal, root cause, files touched, decisions made, commands run, failing or passing tests, and the next concrete step."},
			{Kind: types.Code, Content: "  Good compact handoff:\n  - Goal: fix parser panic on empty input\n  - Evidence: TestParseEmpty panics in parser.go:42\n  - Changed: parser.go guard before token lookup\n  - Verified: go test ./parser passes\n  - Next: run full go test ./... and open PR"},
			{Kind: types.Paragraph, Content: "Compaction is not the same as forgetting. The agent should compact after a major phase, before context gets tight, or before handing work to another session. It should not compact away unresolved evidence during active debugging."},
			{Kind: types.Heading, Content: "Long-term Memory"},
			{Kind: types.Paragraph, Content: "Long-term memory stores durable facts that should survive beyond the current task. Good candidates include stable user preferences, repository conventions, repeated commands, and architectural decisions. Temporary stack traces, current branch names, or half-formed guesses should stay out of memory."},
			{Kind: types.Code, Content: "  Save to memory:\n  - This repo uses table-driven tests\n  - Run make test-unit before committing\n  - User prefers concise final summaries\n\n  Do not save:\n  - Current failing line number\n  - Today's temporary branch\n  - A hypothesis not yet verified"},
			{Kind: types.Heading, Content: "Where Memory Lives"},
			{Kind: types.Paragraph, Content: "\"Long-term memory\" is not one thing. It is a storage choice, and different kinds of facts want different stores. Most real agents combine several rather than picking one."},
			{Kind: types.Code, Content: "  Plain files (Markdown notes, config)\n    Good: human-readable, diffable, version-controlled\n    Weak: no query language, must be read whole\n\n  Key-value store\n    Good: fast exact lookup by a known key\n    Weak: you must already know what to ask for\n\n  Relational / structured store\n    Good: filters, joins, and reporting over facts\n    Weak: needs a schema up front\n\n  Vector store (embeddings)\n    Good: \"find things that mean something similar\"\n    Weak: fuzzy, needs re-indexing, can return confident noise\n\n  Knowledge graph\n    Good: relationships between entities and decisions\n    Weak: expensive to build and keep accurate"},
			{Kind: types.Paragraph, Content: "A useful default for a coding agent: version-controlled files for conventions the team should also read, a key-value or structured store for per-user preferences, and a vector store only once the corpus is too large to load or grep directly."},
			{Kind: types.Heading, Content: "Retrieval and RAG"},
			{Kind: types.Paragraph, Content: "Storage is only half the problem. Once memory outgrows the context window, the agent has to fetch the relevant slice at the right moment. That is retrieval, and Retrieval-Augmented Generation (RAG) is the common name for the pattern: retrieve relevant material first, then generate an answer grounded in what was retrieved."},
			{Kind: types.Code, Content: "  RAG pipeline\n\n  1. Chunk     split documents into retrievable units\n  2. Embed     turn each chunk into a vector\n  3. Index     store vectors for nearest-neighbor search\n  4. Retrieve  embed the query, pull the top-k chunks\n  5. Re-rank   reorder or filter the candidates\n  6. Generate  answer using only the retrieved context"},
			{Kind: types.Paragraph, Content: "RAG keeps the context window small while letting the agent draw on a corpus far larger than the window. It also gives you citations: the agent can point at which chunk a claim came from, which makes wrong answers debuggable."},
			{Kind: types.Bullet, Content: "Retrieval quality caps answer quality - a bad top-k makes a confident wrong answer\nChunk boundaries matter - a chunk that splits a function or an argument retrieves badly\nStale indexes lie - re-index when the underlying files change\nSemantic search is not always the answer - exact grep or symbol lookup often beats embeddings in a codebase\nRetrieved text is data, not instructions - never let a retrieved chunk redirect the agent"},
			{Kind: types.Heading, Content: "Agent State Design"},
			{Kind: types.Bullet, Content: "Active context - detailed working set for the current task\nCompacted state - concise handoff that preserves progress\nLong-term memory - durable facts reused across sessions\nExternal artifacts - commits, test logs, issues, and PRs that outlive chat history"},
			{Kind: types.Heading, Content: "Practical Rule"},
			{Kind: types.Paragraph, Content: "Before saving anything, ask: will this still be true and useful next week? Before compacting, ask: could another agent continue the task from this summary without re-reading everything?"},
			{Kind: types.Callout, Content: "Reliable agents separate short-lived working context from durable memory. That separation prevents stale assumptions while still preserving useful experience."},
			{Kind: types.Heading, Content: "Further Reading"},
			{Kind: types.Callout, Content: "Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks (Lewis et al., 2020) — https://arxiv.org/abs/2005.11401"},
			{Kind: types.Callout, Content: "Dense Passage Retrieval for Open-Domain Question Answering (Karpukhin et al., 2020) — https://arxiv.org/abs/2004.04906"},
			{Kind: types.Callout, Content: "Lost in the Middle: How Language Models Use Long Contexts (Liu et al., 2023) — https://arxiv.org/abs/2307.03172"},
			{Kind: types.Callout, Content: "Vector database — https://en.wikipedia.org/wiki/Vector_database"},
		},
		Questions: []types.QuizQuestion{
			{
				Kind:        types.MultipleChoice,
				Prompt:      "What belongs in a compaction summary?",
				Choices:     []string{"Every line of every file read", "Goal, evidence, decisions, commands, and next steps", "Only the final answer", "Only user preferences"},
				CorrectIdx:  1,
				Explanation: "A compact summary should preserve the state needed to continue work without replaying the whole conversation.",
			},
			{
				Kind:        types.FillBlank,
				Prompt:      "Stable repo conventions should be stored in long-term ____.",
				Answer:      "memory",
				Explanation: "Long-term memory is for durable information that remains useful across sessions.",
			},
			{
				Kind:        types.MultipleChoice,
				Prompt:      "Which item should NOT be saved as long-term memory?",
				Choices:     []string{"The repo uses go vet before PRs", "The user prefers concise final answers", "The current stack trace from one debugging run", "The team uses table-driven tests"},
				CorrectIdx:  2,
				Explanation: "A stack trace from one debugging run is task-local evidence, not durable memory.",
			},
			{
				Kind:        types.MultipleChoice,
				Prompt:      "A team convention should be readable and reviewable by humans as well as the agent. Which store fits best?",
				Choices:     []string{"A vector store of embeddings", "A version-controlled Markdown file in the repo", "An in-memory cache for this session", "A knowledge graph of entities"},
				CorrectIdx:  1,
				Explanation: "Conventions the team also reads belong in version-controlled files, where they are diffable and reviewable. Vector stores are for corpora too large to read directly.",
			},
			{
				Kind:        types.MultipleChoice,
				Prompt:      "What is the point of RAG in an agent?",
				Choices:     []string{"It removes the need for a context window", "It retrains the model on your repository", "It retrieves the relevant slice of a large corpus and grounds the answer in it", "It compresses the conversation into a summary"},
				CorrectIdx:  2,
				Explanation: "RAG retrieves relevant material first and generates from it, so the agent can draw on a corpus much larger than the context window. Compressing the conversation is compaction, not retrieval.",
			},
		},
	})
}
