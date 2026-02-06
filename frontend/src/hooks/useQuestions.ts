import { useQuery } from '@tanstack/react-query';
import { fetchQuestions } from '../services/problemService';
import { type Question } from '../types/question';

export function useQuestions() {
    return useQuery<Question[], Error>({
        queryKey: ['questions'],
        queryFn: fetchQuestions,
    });
}