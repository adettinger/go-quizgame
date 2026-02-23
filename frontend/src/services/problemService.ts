import { ProblemType, type Problem } from '../types/problem';
import type { Question } from '../types/question';
import type { QuestionSubmission, QuizSubmission } from '../types/requests';
import type { StartQuizResponse, SubmitQuizResponse } from '../types/responses';

const API_URL = 'http://localhost:8080';

export async function fetchProblems(): Promise<Problem[]> {
    const response = await fetch(`${API_URL}/problem`);

    if (!response.ok) {
        throw new Error(`Error fetching problems: ${response.status}`);
    }

    return response.json();
}

export async function startQuiz(): Promise<StartQuizResponse> {
    const response = await fetch(`${API_URL}/quiz/start`);

    if (!response.ok) {
        throw new Error(`Error fetching questions: ${response.status}`);
    }

    return response.json();
}

export async function fetchQuestions(): Promise<Question[]> {
    const response = await fetch(`${API_URL}/quiz/questions`);

    if (!response.ok) {
        throw new Error(`Error fetching questions: ${response.status}`);
    }

    return response.json();
}

export async function fetchProblemById(id: string): Promise<Problem> {
    const response = await fetch(`${API_URL}/problem/${id}`);

    if (!response.ok) {
        throw new Error(`Error fetching problems: ${response.status}`);
    }

    return response.json();
}

// TODO: Move interfaces
interface ProblemFormData {
    Type: ProblemType,
    Question: string;
    Answer: string;
    Choices: string[];
}

export async function createProblem(data: ProblemFormData): Promise<Problem> {
    const response = await fetch(`${API_URL}/problem`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        throw new Error(`Error creating problem: ${response.status}`);
    }

    return response.json();
};

export async function submitQuiz(data: QuizSubmission): Promise<SubmitQuizResponse> {
    const response = await fetch(`${API_URL}/quiz/submit`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error(`Error submitting quiz: ${response.status}`);
    }

    return response.json();
}

export async function deleteProblemById(id: string): Promise<string> {
    const response = await fetch(`${API_URL}/problem/${id}`, {
        method: 'DELETE',
    });

    if (!response.ok) {
        throw new Error(`Error deleting problems: ${response.status}`);
    }

    return id;
}