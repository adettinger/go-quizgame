import { CreateProblemForm } from "../components/CreateProblemForm";

export function CreateProblemPage() {
    return (
        <div className="container mx-auto p-4">
            <h1 className="text-3xl font-bold mb-6">Create Problem</h1>
            <CreateProblemForm />
        </div>
    );
}