import welcomImage from '../assets/WelcomeNoneOfYourQuizness.jpg';

function Welcome() {
    return (
        <div className="text-center mt-5">
            <h1>Welcome to None of your Quizness!</h1>
            <img
                src={welcomImage}
                alt={'Welcome'}
                style={{ width: 'auto' }}
                loading="lazy"
            />
        </div>
    );
}

export default Welcome;