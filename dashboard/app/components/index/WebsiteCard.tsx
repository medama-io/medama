import { Text, UnstyledButton } from '@mantine/core';
import { Link } from '@remix-run/react';

import classes from './WebsiteCard.module.css';

interface WebsiteCardProps {
	website: {
		hostname: string;
		summary?: {
			visitors: number;
		};
	};
}

export const WebsiteCard = ({ website }: WebsiteCardProps) => {
	return (
		<UnstyledButton
			key={website.hostname}
			p={16}
			component={Link}
			to={`/${website.hostname}`}
			prefetch="intent"
			className={classes.card}
			role="link"
			tabIndex={0}
			aria-label={`Visit ${website.hostname}`}
			aria-describedby={`${website.hostname}-visitors`}
		>
			<Text truncate="end">{website.hostname}</Text>
			<Text
				component="span"
				size="xs"
				c="gray"
				id={`${website.hostname}-visitors`}
			>
				{website.summary?.visitors ?? 'N/A'} visitors in the last 24 hours
			</Text>
		</UnstyledButton>
	);
};
