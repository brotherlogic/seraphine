# Seraphine

Seraphine is a workflow and system for running code - it is planning and review centric with an eye on putting humans involvement at the start and end of code changes - i.e. in the feature development and review stages.

## Workflow

The workflow consists of three steps - we start from a broadly defined github issue which outlines a new features. We then refine the feature and tease out the requirements, plan out some work, breaking it down into managable features and then review the code that implements the feature.

## Init

To init a project we run:

seraphine init

in the project root. This adds .seraphine directory to the project and gives seraphine the necessary permissions to push changes and manage issues etc.

## Feature planning

To build out a new feature we either create an issue with a broadly defined statement.