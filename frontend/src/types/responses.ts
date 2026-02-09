import type { Question } from "./question";

export interface SubmitQuizResponse {
    Score: number;
    Answers: QuestionResponse[];
}

export interface QuestionResponse {
    Id: string;
    Answer: string;
    Correct: boolean;
}

export interface StartQuizResponse {
    SessionId: string;
    Timeout: EpochTimeStamp;
    Questions: Question[];
}