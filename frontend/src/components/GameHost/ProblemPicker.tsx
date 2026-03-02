import { Flex, Text } from "@radix-ui/themes";
import type { Problem } from "../../types/problem";
import React, { useEffect, useState } from "react";
import { useProblems } from "../../hooks/useProblems";
import { ProblemTable } from "../ProblemTable";

export interface ProblemPickerProps {
    setQuestionIds: (questionIds: string[]) => void;
};

export const ProblemPicker = React.memo(function ({ setQuestionIds }: ProblemPickerProps) {
    const { data, isLoading, isError, error } = useProblems();

    const [availableQuestions, setAvailableQuestions] = useState<Problem[]>([]);
    const [selectedQuestions, setSelectedQuestions] = useState<Problem[]>([]);

    useEffect(() => {
        if (!!data && data.length > 0) {
            setAvailableQuestions(data);
        }
    }, [data]);

    useEffect(() => {
        setQuestionIds(selectedQuestions.map(problem => problem.Id))
    }, [selectedQuestions])

    const selectQuestion = (problem: Problem) => {
        setSelectedQuestions(prev => [...prev, problem]);
        setAvailableQuestions(prev => prev.filter(question => question.Id !== problem.Id));
    };

    const deselectQuestion = (problem: Problem) => {
        setSelectedQuestions(prev => prev.filter(question => question.Id !== problem.Id));
        setAvailableQuestions(prev => [...prev, problem]);
    };

    if (isLoading) {
        return <div className="text-center p-4">Loading problems...</div>;
    }

    if (isError && !!error) {
        return (
            <div className="text-center p-4 text-red-600">
                Error loading problems: {error.message}
            </div>
        );
    }

    if (availableQuestions.length === 0 && selectedQuestions.length === 0) {
        return <div className="text-center p-4">No problems found.</div>;
    }

    return (
        <Flex direction="column" align="center" justify="center" style={{ width: '100%' }}>
            <h3>Selected Questions</h3>
            {selectedQuestions.length === 0 ?
                <Text>No Questions Selected</Text>
                :
                <ProblemTable
                    Problems={selectedQuestions}
                    ShowIds={false}
                    OnRemove={deselectQuestion}

                />
            }

            <h3>Available Questions</h3>
            {availableQuestions.length === 0 ?
                <Text>No available questions remaining</Text>
                :
                <ProblemTable
                    Problems={availableQuestions}
                    ShowIds={false}
                    OnAdd={selectQuestion}

                />
            }
        </Flex>
    );
});