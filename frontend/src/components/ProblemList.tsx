import { Link, useNavigate } from 'react-router-dom';
import { useProblems } from '../hooks/useProblems';
import { Button, DropdownMenu, Flex, IconButton, Table } from "@radix-ui/themes"
import { Pencil1Icon, TrashIcon } from '@radix-ui/react-icons';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteProblemById } from '../services/problemService';
import { useToast } from './Toast/ToastContext';
import { ProblemType } from '../types/problem';
import { ProblemTable } from './ProblemTable';

export function ProblemList() {
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const { showToast } = useToast();
    const { data, isLoading, isError, error } = useProblems();

    const deleteMutation = useMutation({
        mutationFn: deleteProblemById,
        onSuccess: (id: string) => {
            queryClient.invalidateQueries({ queryKey: ['problems'] });
            queryClient.invalidateQueries({ queryKey: ['problem', id] });
            showToast('success', "Success", `Deleted problem ${id}`);
            console.log(`Deleted problem ${id}`);
        },
        onError: (error) => {
            console.log('Failed to delete problem', error)
        },
    })

    const handleDelete = (id: string) => {
        if (confirm('Are you sure you want to delete this problem?')) {
            deleteMutation.mutate(id);
        }
    }

    if (isLoading) {
        return <div className="text-center p-4">Loading problems...</div>;
    }

    if (isError) {
        return (
            <div className="text-center p-4 text-red-600">
                Error loading problems: {error.message}
            </div>
        );
    }

    if (!data || data.length === 0) {
        return <div className="text-center p-4">No problems found.</div>;
    }

    const handleEdit = (id: string) => {
        navigate(`/problem/edit/${id}`);
    };


    return (
        <Flex align="center" justify="center">
            <ProblemTable
                Problems={data}
                IsLoading={isLoading}
                IsError={isError}
                Error={error}
                ShowIds={true}
                OnEdit={handleEdit}
                OnDelete={handleDelete}
                DisableActions={deleteMutation.isPending}
            />
        </Flex>
    );
}