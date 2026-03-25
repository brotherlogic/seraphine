# Triage Agent

The Triage agent is a mid level agent that scans the github repository, looks in the /features directory and finds any .md file that has no status line.

A status line is one that has "status: " and then one of CONSIDERING, PROPOSED, PLANNED, IMPLEMENTING, COMPLETE following the colon.

If the agent finds such a feature it:
  1. Validates that the feature makes sense - adding follow up questions in the body of the file (in the format TRIAGE: xxxx) where xxx is the question.
  1. Once the triage comments are addressed, it adds a
     status line in the top level of the Markdown file of the format STATUS: PROPOSED.
  1. It then follows the standard .agent/workflows/finish.md script

The validation should use the most advanced gemini model at the time of running, and trigger the prompt:

"Are there any clarifications you'd make to this proposal, or is it ready to code"

And the clarifications are written into the document and pushed to the same branch.