# 🕵️ The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the agent must execute a deep research process to identify and evaluate potential solutions to the problem space.

## 🔄 Workflow Lifecycle

1. **Search**: Evaluate the problem space and search for potential solutions or options. Continue researching until at least 3 viable options are found. You should utilize the `deep-research` skill if available.
2. **Synthesize**: Organically synthesize the findings, outlining the pros and cons of the identified options.
3. **Comment**: Post the synthesized summary as a comment on the parent issue to record the findings.
4. **Remove Label**: Remove the `seraphine-needs-deep-research` label from the issue.
5. **Create Sub-issue**: Create a new sub-issue titled `[Requirements] <Original Issue Title>` and label it with `seraphine-needs-requirements` to initiate the next phase (requirements gathering).

## ⚠️ Error States & Handling

- **No Viable Options Found**: If exhaustive research yields no viable options, post a comment explaining the dead-end and the avenues explored. Remove the `seraphine-needs-deep-research` label. Do NOT create a Requirements sub-issue. (Optional: Apply a `seraphine-blocked` or similar label if appropriate, and stop execution).
- **Ambiguous Problem Space**: If the issue description is too vague to proceed, comment on the issue requesting clarification from the author. Remove the `seraphine-needs-deep-research` label, and wait for further instructions. Do NOT create a Requirements sub-issue.
