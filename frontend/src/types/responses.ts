export interface SubmitQuizResponse {
    Score: number;
    Answers: QuestionResponse[];
}

export interface QuestionResponse {
    Id: string;
    Answer: string;
    Correct: boolean;
}