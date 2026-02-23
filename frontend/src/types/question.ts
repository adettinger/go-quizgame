import type { ProblemType } from "./problem";

export interface Question {
    Id: string;
    Type: ProblemType;
    Question: string;
    Choices: string[];
}