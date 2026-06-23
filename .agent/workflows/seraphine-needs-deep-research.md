# 🔍 The `seraphine-needs-deep-research` Label Workflow

When a parent issue is labeled with `seraphine-needs-deep-research`, the AI assistant is triggered to execute a deep research phase to identify options and synthesize findings before transitioning to requirements gathering.

## 🔄 Workflow Sequence

1. **Search**: Invoke the `deep-research` skill to explore the problem space. Evaluate the problem space until at least the top 3 viable options are identified.
2. **Synthesize**: Organically summarize the pros and cons of these options (no strict formatting template required).
3. **Comment**: Post the synthesized summary as a comment to the parent issue.
4. **Remove Label**: Programmatically remove the `seraphine-needs-deep-research` label from the parent issue using the `gh` CLI.
5. **Create Sub-issue**: Programmatically create a new child `[Requirements]` sub-issue using the GitHub API (or `gh` CLI) and label it with `seraphine-needs-requirements` to transition to the requirements gathering phase.

## ⚠️ Error States & Handling

* **No Viable Options Found**: If the research yields no viable options, document the constraints and reasons why no options exist. Post this as a comment to the parent issue, leave the label intact, and explicitly flag the user for input. Do not create the child requirements sub-issue.
* **Ambiguous Problem Space**: If the problem space is too broad or ambiguous to identify 3 distinct options, summarize the ambiguity and the areas that need clarification. Post this summary as a comment, request user clarification, and do not remove the label or create the child sub-issue until clarification is provided.
