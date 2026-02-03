import { ProblemList } from '../components/ProblemList';

export function ProblemsPage() {
    return (
        <div className="container mx-auto p-4">
            <h1 className="text-3xl font-bold mb-6">Problem List</h1>
            <ProblemList />
        </div>
    );
}