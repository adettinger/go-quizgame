import { Button, Flex, Switch, TextField, Tooltip } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useState } from "react";
import { ProblemPicker } from "./ProblemPicker";
import { QuestionMarkCircledIcon } from "@radix-ui/react-icons";


export function GameOptions() {
    const [questionTimeLimit, setQuestionTimeLimit] = useState("0");
    const [isHostAPlayer, setIsHostAPlayer] = useState(false);
    const [selectedQuestionIds, setSelectedQuestionIds] = useState<string[]>([]);


    const areAllFieldsValid = (): boolean => {
        if (
            parseInt(questionTimeLimit, 10) >= 0 &&
            selectedQuestionIds.length > 0
        ) {
            return true
        }
        return false;
    };

    const handleSubmit = () => {
        // TODO: Actually start a game
        console.log("Starting a game", { "Host is a player": isHostAPlayer, "question time limit": questionTimeLimit, "selected questions": selectedQuestionIds });
    };

    return (
        <Form.Root onSubmit={handleSubmit}>
            <Flex direction="column" justify={"center"} align={"center"} gap="3">
                {/* Time limit on questions */}
                <Form.Field name="timeLimit">
                    <Flex direction={"row"} gap="1">
                        <Form.Label>Time limit per question (seconds):</Form.Label>
                        <Tooltip content={"0 for no time limit"}>
                            <QuestionMarkCircledIcon />
                        </Tooltip>
                        <Form.Control asChild>
                            <TextField.Root
                                value={questionTimeLimit}
                                onChange={(event) => { setQuestionTimeLimit(event.target.value.replace(/[^0-9\s]/g, '')) }}
                            />
                        </Form.Control>
                    </Flex>
                </Form.Field>
                {/* is host a player */}
                <Form.Field name="isHostAPlayer">
                    <Flex direction="row" gap="3">
                        <Form.Label>Host is a player:</Form.Label>
                        <Form.Control asChild>
                            <Switch checked={isHostAPlayer} onCheckedChange={setIsHostAPlayer} />
                        </Form.Control>
                    </Flex>
                </Form.Field>

                {/* question list */}
                <Form.Field name="selectedQuestions">
                    <Form.Control asChild>
                        <ProblemPicker
                            setQuestionIds={setSelectedQuestionIds}
                        />
                    </Form.Control>
                </Form.Field>

                <Form.Submit asChild>
                    <Button disabled={!areAllFieldsValid()}>Start a game!</Button>
                </Form.Submit>
            </Flex>
        </Form.Root>
    );
};