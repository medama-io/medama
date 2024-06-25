import { Group, Text, UnstyledButton } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';

import React, { useMemo } from 'react';
import { formatCount, formatDuration } from './formatter';
import classes from './StatsDisplay.module.css';

interface StatsItemProps {
	label: string;
	count?: number;
	percentage?: number;
	tab: string;
}

const FILTER_MAP: Record<string, string> = {
	Pages: 'path',
	Time: 'path',
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

const StatsItem = ({
	label,
	count = 0,
	percentage = 0,
	tab,
}: StatsItemProps) => {
	const [searchParams, setSearchParams] = useSearchParams();

	const formattedValue = useMemo(
		() => (tab === 'Time' ? formatDuration(count) : formatCount(count)),
		[tab, count],
	);

	const handleFilter = () => {
		const params = new URLSearchParams(searchParams);
		const filter = FILTER_MAP[tab] ?? 'path';
		params.append(`${filter}[eq]`, label);
		setSearchParams(params, { preventScrollReset: true });
	};

	return (
		<UnstyledButton
			className={classes['stat-item']}
			onClick={handleFilter}
			aria-label={`Filter by ${label}`}
			disabled={tab === 'Time'}
		>
			<Group justify="space-between" pb={6}>
				<Text fz={14} truncate>
					{label}
				</Text>
				<Text fw={600} fz={14}>
					{formattedValue}
				</Text>
			</Group>
			<div
				className={classes.bar}
				style={{ width: `${Math.min(percentage * 100, 100)}%` }}
				aria-hidden="true"
			/>
		</UnstyledButton>
	);
};

const StatsItemMemo = React.memo(StatsItem);

export { StatsItemMemo as StatsItem };
