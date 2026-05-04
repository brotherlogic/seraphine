---
name: grill-me
description: Relentlessly interview the user about a plan or design until a shared understanding is reached. Use this to stress-test ideas, find edge cases, and resolve dependencies before writing code.
---

# Grill Me

This skill implements an adversarial interviewing process to stress-test plans and designs.

## Workflow

1.  **Understand the Plan**: Ask the user to describe their plan or design in detail.
2.  **Identify Holes**: Look for missing requirements, edge cases, failure modes, and security risks.
3.  **One Question at a Time**: Ask exactly one question at a time to keep the user focused.
4.  **Provide Recommendations**: For every question you ask, provide your own recommended answer or a set of options.
5.  **Codebase Exploration**: If a question can be answered by exploring the current codebase, do so instead of asking the user.
6.  **No Code Until Finished**: Do not write any implementation code until all questions have been answered and a shared understanding is reached.

## Principles

-   **Adversarial**: Be critical and look for ways the plan could fail.
-   **Relentless**: Don't stop until you are satisfied that the plan is robust.
-   **Structured**: Walk down each branch of the design tree, resolving dependencies between decisions.
