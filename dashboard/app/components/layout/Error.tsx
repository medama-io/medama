import { Anchor, Container, Text, Title } from '@mantine/core';
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

const ErrorPage = ({ label, title, description }: ErrorPageProps) => {
	return (
		<Container className={classes.root}>
			<div className={classes.label}>{label}</div>
			<Title className={classes.title}>{title}</Title>
			<Text c="dimmed" size="lg" ta="center" className={classes.description}>
				{description}
			</Text>
		</Container>
	);
};

const BadRequestError = () => (
	<ErrorPage
		label="400"
		title="Your request is invalid."
		description={
			<>
				Please check the URL and try again. If you think this is an error,
				please{' '}
				<Anchor
					component={Link}
					to="https://github.com/medama-io/medama/issues"
					className={classes.anchor}
				>
					report this issue
				</Anchor>{' '}
				to the developers.
			</>
		}
	/>
);

const ForbiddenError = () => {
	// Check if hostname is demo.medama.io or medama.fly.dev to display a different message.
	const hostname = window.location.hostname;
	const isDemo = hostname === 'demo.medama.io' || hostname === 'medama.fly.dev';

	const description = isDemo ? (
		<>
			You are currently in demo mode. You can't access this page or perform this
			action.
		</>
	) : (
		<>
			You don't have permission to view this page or perform this action. Please
			contact your administrator if you believe this is an error.
		</>
	);

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
				<Anchor component={Link} to="/" className={classes.anchor}>
					home page
				</Anchor>{' '}
				or browse the docs to find what you need.
			</>
		}
	/>
);

const ConflictError = () => (
	<ErrorPage
		label="409"
		title="Conflict detected."
		description={
			<>
				There was a conflict while processing your request. Did you try to
				create something that already exists?
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
				<Anchor
					component={Link}
					to="https://github.com/medama-io/medama/issues"
					className={classes.anchor}
				>
					report this issue
				</Anchor>{' '}
				to the developers.
				<br />
				<Text component="span" c="red">
					{error ? `Error: ${error}` : ''}
				</Text>
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
};
