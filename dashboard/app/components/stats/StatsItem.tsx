import { Group, Text, UnstyledButton } from '@mantine/core';
import { useHover } from '@mantine/hooks';
import React, { useMemo } from 'react';
import isFQDN from 'validator/lib/isFQDN';

import { IconExternal } from '@/components/icons/external';
import { useFilter } from '@/hooks/use-filter';

import { formatCount, formatDuration } from './formatter';
import type { Filter } from './types';

import classes from './StatsItem.module.css';

interface StatsItemProps {
	label: string;
	count?: number;
	percentage?: number;
	tab: string;
}

const FILTER_MAP: Record<string, Filter> = {
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
	const { addFilter, isFilterActiveEq } = useFilter();
	const showLocales = isFilterActiveEq('language');
	const { hovered, ref } = useHover<HTMLButtonElement>();

	const formattedValue = useMemo(
		() => (tab === 'Time' ? formatDuration(count) : formatCount(count)),
		[tab, count],
	);

	const handleFilter = () => {
		const filter = FILTER_MAP[tab] ?? 'path';
		if (filter !== 'language' || !showLocales) {
			addFilter(filter, 'eq', label);
		}
	};

	return (
		<UnstyledButton
			className={classes.item}
			onClick={handleFilter}
			aria-label={`Filter by ${label}`}
			ref={ref}
		>
			<Group justify="space-between" pb={6} wrap="nowrap">
				<Group gap="xs">
					<Text fz={14} truncate style={{ userSelect: 'text' }}>
						{label}
					</Text>
					{tab === 'Referrers' && (
						<UnstyledButton
							className={classes.external}
							component="a"
							href={`https://${label}`}
							target="_blank"
							rel="noreferrer noopener"
							data-hidden={!isFQDN(label)}
							onClick={(event) => event.stopPropagation()}
						>
							<IconExternal />
						</UnstyledButton>
					)}
				</Group>
				<Group align="center" gap="xs" wrap="nowrap">
					<Text
						component="span"
						fz={12}
						c="gray"
						mr={4}
						data-active={hovered ? 'true' : undefined}
						visibleFrom="xs"
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
