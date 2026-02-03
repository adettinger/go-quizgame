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