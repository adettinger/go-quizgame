import { useQuery } from '@tanstack/react-query';
import { fetchProblemById } from '../services/problemService';
import { type Problem } from '../types/problem';

export function useProblem(id: string) {
    return useQuery<Problem, Error>({
        queryKey: ['problem', id],
        queryFn: () => fetchProblemById(id),
        enabled: !!id,
    });
}