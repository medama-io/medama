import { Group, Text, UnstyledButton } from '@mantine/core';
import { useSearchParams } from '@remix-run/react';
import { useHover } from '@mantine/hooks';

import React, { useMemo } from 'react';
import { formatCount, formatDuration } from './formatter';

import classes from './StatsItem.module.css';

interface StatsItemProps {
	label: string;
	count?: number;
	percentage?: number;
	tab: string;
}

const FILTER_MAP: Record<string, string> = {
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
	const { hovered, ref } = useHover<HTMLButtonElement>();

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
			className={classes.item}
			onClick={handleFilter}
			aria-label={`Filter by ${label}`}
			ref={ref}
		>
			<Group justify="space-between" pb={6}>
				<Text fz={14} truncate>
					{label}
				</Text>
				<Group align="center" gap="xs">
					<Text
						component="span"
						fz={12}
						c="gray"
						mr={4}
						data-active={hovered ? 'true' : undefined}
					>
						{(percentage * 100).toFixed(1)}%
					</Text>
					<Text fw={600} fz={14}>
						{formattedValue}
					</Text>
				</Group>
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
