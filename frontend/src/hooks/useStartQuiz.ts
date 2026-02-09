import { useQuery } from '@tanstack/react-query';
import { startQuiz } from '../services/problemService';
import type { StartQuizResponse } from '../types/responses';

export function useStartQuiz() {
    return useQuery<StartQuizResponse, Error>({
        queryKey: ['startQuiz'],
        queryFn: startQuiz,
    });
}