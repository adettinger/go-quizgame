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

export function getEnumKeyByValue<T extends { [index: string]: string }>(enumObj: T, value: string): keyof T | undefined {
    return Object.keys(enumObj).find(key => enumObj[key] === value) as keyof T | undefined;
}

export const MAX_PROBLEM_CHOICES = 4;