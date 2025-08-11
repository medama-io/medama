import { Link } from '@remix-run/react';
import type { ReactNode } from 'react';

import classes from './Error.module.css';

interface ErrorPageProps {
	label: string;
	title: string;
	description: string | ReactNode;
}

interface InternalServerErrorProps {
	error?: string;
}

interface ErrorProps {
	message?: string;
}

class StatusError extends Error {
	status: number;

	constructor(status: number, message: string) {
		super(message);
		this.status = status;
	}
}

const isStatusError = (error: unknown): error is StatusError => {
	return error instanceof StatusError;
};

const ErrorPage = ({ label, title, description }: ErrorPageProps) => {
	return (
		<div className={classes.root}>
			<div className={classes.label}>{label}</div>
			<h3 className={classes.title}>{title}</h3>
			<p className={classes.description}>{description}</p>
		</div>
	);
};

const BadRequestError = ({ message }: ErrorProps) => (
	<ErrorPage
		label="400"
		title="Your request is invalid."
		description={
			<>
				Please check the URL and try again. If you think this is an error,
				please{' '}
				<a
					className={classes.anchor}
					href="https://github.com/medama-io/medama/issues"
					target="_blank"
					rel="noreferrer"
				>
					report this issue
				</a>{' '}
				to the developers.
				{message && (
					<div>
						<br />
						<span className={classes.error}>
							{message ? `Error: ${message}` : ''}
						</span>
					</div>
				)}
			</>
		}
	/>
);

const ForbiddenError = () => {
	// Check if hostname is demo.medama.io or medama.fly.dev to display a different message.
	const hostname = window.location.hostname;
	const isDemo = hostname === 'demo.medama.io' || hostname === 'medama.fly.dev';

	const description = isDemo
		? "You are currently in demo mode. You can't access this page or perform this action."
		: "You don't have permission to view this page or perform this action. Please contact your administrator if you believe this is an error.";

	return (
		<ErrorPage
			label="403"
			title="You are not allowed to access this page."
			description={description}
		/>
	);
};

const NotFoundError = () => (
	<ErrorPage
		label="404"
		title="You have found a secret place."
		description={
			<>
				The page you're looking for isn't here. Feel free to return to the{' '}
				<Link to="/" className={classes.anchor}>
					home page
				</Link>{' '}
				or browse the docs to find what you need.
			</>
		}
	/>
);

const ConflictError = ({ message }: ErrorProps) => (
	<ErrorPage
		label="409"
		title="Conflict detected."
		description={
			<>
				There was a conflict while processing your request. Did you try to
				create something that already exists?
				{message && (
					<div>
						<br />
						<span className={classes.error}>
							{message ? `Error: ${message}` : ''}
						</span>
					</div>
				)}
			</>
		}
	/>
);

const InternalServerError = ({ error }: InternalServerErrorProps) => (
	<ErrorPage
		label="500"
		title="Something went wrong."
		description={
			<>
				We encountered an unexpected error while processing your request. Please{' '}
				<a
					className={classes.anchor}
					href="https://github.com/medama-io/medama/issues"
					target="_blank"
					rel="noreferrer"
				>
					report this issue
				</a>{' '}
				to the developers.
				<div>
					<br />
					<span className={classes.error}>
						{error ? `Error: ${error}` : ''}
					</span>
				</div>
			</>
		}
	/>
);

export {
	BadRequestError,
	ConflictError,
	ForbiddenError,
	InternalServerError,
	NotFoundError,
	StatusError,
	isStatusError,
};
