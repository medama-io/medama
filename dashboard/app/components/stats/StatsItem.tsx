import React, { useMemo } from 'react';
import isFQDN from 'validator/lib/isFQDN';

import { IconExternal } from '@/components/icons/external';
import { Group } from '@/components/layout/Flex';
import { useFilter } from '@/hooks/use-filter';
import { useHover } from '@/hooks/use-hover';

import { formatCount, formatDuration } from './formatter';
import type { Filter } from './types';

import classes from './StatsItem.module.css';

interface StatsItemProps {
	label: string;
	count?: number;
	percentage?: number;
	tab: string;
	bar?: boolean;
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
	PropName: 'prop_name',
	PropValue: 'prop_value',
};

const StatsItem = ({
	label,
	count = 0,
	percentage,
	tab,
	bar = true,
}: StatsItemProps) => {
	const { addFilter, isFilterActiveEq } = useFilter();
	const { hovered, ref } = useHover<HTMLButtonElement>();

	const formattedValue = useMemo(
		() => (tab === 'Time' ? formatDuration(count) : formatCount(count)),
		[tab, count],
	);

	const handleFilter = () => {
		let key = tab;
		// Properties tab has two types of filters
		if (tab === 'properties') {
			key = isFilterActiveEq('prop_name') ? 'PropValue' : 'PropName';
		}

		const filter = FILTER_MAP[key] ?? 'path';
		addFilter(filter, 'eq', label);
	};

	// If percentage is not defined, don't show the percentage bar
	if (percentage === undefined) {
		bar = false;
		percentage = 0;
	}

	return (
		<button
			className={classes.item}
			type="button"
			onClick={handleFilter}
			aria-label={`Filter by ${label}`}
			ref={ref}
		>
			<Group style={{ paddingBottom: bar ? 6 : 0, flexWrap: 'nowrap' }}>
				<Group style={{ overflow: 'hidden', gap: 8 }}>
					<p className={classes.label}>{label}</p>
					{tab === 'Referrers' && (
						<a
							className={classes.external}
							aria-label={`Visit ${label}`}
							href={`https://${label}`}
							target="_blank"
							rel="noreferrer noopener"
							data-hover={hovered ? 'true' : undefined}
							data-hidden={!isFQDN(label)}
							onClick={(event) => event.stopPropagation()}
						>
							<IconExternal />
						</a>
					)}
				</Group>
				<Group style={{ gap: 8, flexWrap: 'nowrap' }}>
					{bar && (
						<span
							className={classes.percentage}
							data-active={hovered ? 'true' : undefined}
						>
							{(percentage * 100).toFixed(1)}%
						</span>
					)}
					<p style={{ fontWeight: 600, fontSize: 14 }}>{formattedValue}</p>
				</Group>
			</Group>
			{bar && (
				<div
					className={classes.bar}
					style={{ width: `${Math.min(percentage * 100, 100)}%` }}
					aria-hidden="true"
				/>
			)}
		</button>
	);
};

const StatsItemMemo = React.memo(StatsItem);

export { StatsItemMemo as StatsItem };
