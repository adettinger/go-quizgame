export interface Problem {
    Id: string;
    Type: ProblemType;
    Question: string;
    Choices: string[];
    Answer: string;
}

export enum ProblemType {
    Text = "text",
    Choice = "choice",
}