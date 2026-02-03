import { Flex } from "@radix-ui/themes";
import { ViewProblem } from "../components/ViewProblem/ViewProblem";


export function ViewProblemPage() {
    return (
        <div className="container mx-auto p-4">
            <h1 className="text-3xl font-bold mb-6">View Problem</h1>
            <Flex justify="center">
                <ViewProblem />
            </Flex>
        </div>
    );
}