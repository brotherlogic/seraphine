# 🕵️ The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the agent invokes a dedicated research skill to evaluate the problem space, stop researching once viable options are found, and synthesize the findings.

## 🔄 Workflow Sequence

1. **Search**: Invoke the `deep-research` skill to explore the problem space. The skill will guide you to evaluate options until at least the top 3 viable options are identified.
2. **Synthesize**: Organically synthesize an unbiased summary of the findings, presenting the pros and cons of the discovered options.
3. **Comment**: Post the synthesized summary as a comment directly on the parent issue.
4. **Remove Label**: Programmatically remove the `seraphine-needs-deep-research` label from the issue.
5. **Create Sub-issue**: Transition to the next phase by creating a new `[Requirements]` sub-issue and labeling it with `seraphine-needs-requirements`. Make sure it links to the parent issue.

## ⚠️ Error States

### "No Viable Options Found"
If the research step fails to discover at least 3 viable options after an exhaustive search:
* Clearly document the search avenues attempted.
* State explicitly that sufficient viable options could not be found.
* Summarize whatever options or partial solutions were uncovered.
* Post this summary on the parent issue, remove the label, and create a `[Requirements]` sub-issue tagged with `seraphine-needs-requirements` referencing the limited findings.

### "Ambiguous Problem Space"
If the problem space is too broad or ambiguous to properly research:
* Document the ambiguity encountered (e.g., lack of specific constraints, too many unrelated domains).
* Propose clarifying questions or highlight the missing parameters.
* Post this assessment to the parent issue.
* Remove the `seraphine-needs-deep-research` label and create a `[Requirements]` sub-issue tagged with `seraphine-needs-requirements` to kick off a requirements gathering phase centered on clarifying the ambiguity.
