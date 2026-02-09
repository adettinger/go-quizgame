export interface QuizSubmission {
    SessionId: string;
    Questions: QuestionSubmission[];
}

export interface QuestionSubmission {
    Id: string;
    Answer: string;
}