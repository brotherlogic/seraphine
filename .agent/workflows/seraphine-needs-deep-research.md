# 🕵️ The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the agent must thoroughly investigate the problem space, potential solutions, and external context before progressing.

## 🔄 Workflow Lifecycle

1. **Search**: Gather information from the codebase, web, and external documentation.
2. **Synthesize**: Compile findings into a structured summary of options, tradeoffs, and recommendations.
3. **Comment**: Post the synthesized research findings as a comment on the issue.
4. **Remove Label**: Remove the `seraphine-needs-deep-research` label.
5. **Create Sub-issue**: Add the `seraphine-needs-requirements` label to advance the issue to the next stage.

---

## 📋 Phase Guidelines

### 1. Search
* Perform deep semantic searches across the entire repository.
* Execute targeted web queries for relevant open-source libraries, documentation, or best practices.
* Review any architectural decision records (ADRs) and product requirements documents (PRDs).

### 2. Synthesize
* Organize findings into distinct, viable approaches.
* For each approach, document the pros, cons, implementation complexity, and risks.
* End the synthesis with a clear recommendation based on the current architecture and requirements.

### 3. Comment
* Format the synthesized research as a detailed GitHub comment.
* Ensure markdown elements like tables, code blocks, and bold text are used for readability.

### 4. Remove Label
* Once the research comment is posted successfully, remove the `seraphine-needs-deep-research` label using the `gh` CLI.

### 5. Create Sub-issue / Transition State
* Apply the `seraphine-needs-requirements` label to the current issue to trigger the next phase of processing, as per the issue state transitions.

---

## ⚠️ Error States & Handling

### No Viable Options Found
If research yields no workable solutions:
1. Document the constraints preventing a solution in the issue comment.
2. Tag the author/product owner for clarification.
3. Remove the `seraphine-needs-deep-research` label.
4. Do **not** apply the `seraphine-needs-requirements` label; leave the issue in an open, unlabeled state for human review.

### Ambiguous Problem Space
If the original issue description is too vague to search effectively:
1. Comment detailing what specific information is missing or unclear.
2. Remove the `seraphine-needs-deep-research` label.
3. Do **not** apply the `seraphine-needs-requirements` label.
