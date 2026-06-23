# 🕵️ The `seraphine-needs-deep-research` Label Workflow

When an issue is labeled with `seraphine-needs-deep-research`, the AI assistant is triggered to perform a comprehensive investigation to gather necessary context, explore potential solutions, and synthesize findings before implementation can begin.

## 🔄 Workflow Sequence

1. **Search**
   - Conduct thorough research on the topic, problem space, or requested feature.
   - Investigate the codebase, external documentation, and relevant discussions.
   - Identify potential technical approaches and constraints.

2. **Synthesize**
   - Consolidate the findings from the research phase.
   - Evaluate the pros and cons of the identified approaches.
   - Formulate a clear recommendation or a set of viable options.

3. **Comment**
   - Post a detailed summary of the research and synthesis as a comment on the issue.
   - Ensure the comment includes the identified options, constraints, and recommendations.

4. **Remove Label**
   - Remove the `seraphine-needs-deep-research` label from the issue to indicate the research phase is complete.

5. **Create Sub-issue**
   - Based on the findings, create actionable sub-issues or update the parent issue to proceed to the next stage (e.g., `seraphine-needs-implementation-plan` or `seraphine-needs-requirements`).

---

## ⚠️ Error States & Handling

During the research process, you may encounter specific blockers. Handle them as follows:

### "No Viable Options Found"
- **Trigger:** When research exhausts all avenues and no technically feasible solution exists within current constraints.
- **Action:** 
  1. Document exactly why each explored avenue failed (e.g., API limitations, incompatible dependencies).
  2. Post this documentation as a comment.
  3. Suggest alternative product requirements or relaxed constraints that might make a solution possible.
  4. Remove the `seraphine-needs-deep-research` label and wait for human input.

### "Ambiguous Problem Space"
- **Trigger:** When the original issue is too vaguely defined to effectively research, or multiple contradictory interpretations exist.
- **Action:**
  1. Comment detailing the specific ambiguities encountered.
  2. Ask targeted, clarifying questions to the issue author.
  3. Remove the `seraphine-needs-deep-research` label and apply `seraphine-needs-requirements` to kick it back to the requirements gathering phase.
