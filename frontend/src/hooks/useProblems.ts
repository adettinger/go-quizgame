import { useQuery } from '@tanstack/react-query';
import { fetchProblems } from '../services/problemService';
import { type Problem } from '../types/problem';

export function useProblems() {
    return useQuery<Problem[], Error>({
        queryKey: ['problems'],
        queryFn: fetchProblems,
        // staleTime: Infinity,
        // refetchOnWindowFocus: false,
        // refetchOnMount: false,
    });
}