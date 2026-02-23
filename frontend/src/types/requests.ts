export interface QuizSubmission {
    SessionId: string;
    QuestionSubmissions: QuestionSubmission[];
}

export interface QuestionSubmission {
    QuestionId: string;
    Answer: string;
}