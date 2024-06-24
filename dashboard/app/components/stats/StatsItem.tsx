import { Group, Text, UnstyledButton } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';

import { formatCount, formatDuration } from './formatter';
import classes from './StatsDisplay.module.css';

interface StatsItemProps {
	label: string;
	count?: number;
	percentage?: number;
	tab: string;
}

const filterMap: Record<string, string> = {
	Referrers: 'referrer',
	Sources: 'utm_source',
	Mediums: 'utm_medium',
	Campaigns: 'utm_campaign',
	Browsers: 'browser',
	OS: 'os',
	Devices: 'device',
	Countries: 'country',
	Languages: 'language',
};

const StatsItem = ({ label, count, percentage, tab }: StatsItemProps) => {
	const [searchParams, setSearchParams] = useSearchParams();

	const formattedValue =
		tab === 'Time' ? formatDuration(count ?? 0) : formatCount(count ?? 0);

	const handleFilter = () => {
		if (tab !== 'Time') {
			const params = new URLSearchParams(searchParams);
			const filter = filterMap[tab] ?? 'path';

			params.append(`${filter}[eq]`, label);
			setSearchParams(params, { preventScrollReset: true });
		}
	};

	return (
		<UnstyledButton
			className={classes['stat-item']}
			onClick={handleFilter}
			aria-label={`Filter by ${label}`}
		>
			<Group justify="space-between" pb={6}>
				<Text fz={14}>{label}</Text>
				<Text fw={600} fz={14}>
					{formattedValue}
				</Text>
			</Group>
			<div
				className={classes.bar}
				style={{ width: `${(percentage ?? 0) * 100}%` }}
				aria-hidden="true"
			/>
		</UnstyledButton>
	);
};

export { StatsItem };
