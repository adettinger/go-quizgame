import { type Problem } from '../types/problem';

const API_URL = 'http://localhost:8080';

export async function fetchProblems(): Promise<Problem[]> {
    const response = await fetch(`${API_URL}/problem`);

    if (!response.ok) {
        throw new Error(`Error fetching problems: ${response.status}`);
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

interface ProblemFormData {
    Question: string;
    Answer: string;
}

export async function createProblem(data: ProblemFormData): Promise<Problem> {
    const response = await fetch('http://localhost:8080/problem', {
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

export async function deleteProblemById(id: string): Promise<string> {
    const response = await fetch(`${API_URL}/problem/${id}`, {
        method: 'DELETE',
    });

    if (!response.ok) {
        throw new Error(`Error deleting problems: ${response.status}`);
    }

    return id;
}