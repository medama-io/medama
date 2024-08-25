import * as Tabs from '@radix-ui/react-tabs';
import { useSearchParams } from '@remix-run/react';

import { ButtonLink } from '@/components/Button';
import { useFilter } from '@/hooks/use-filter';

import { Combobox } from './Combobox';
import { StatsItem } from './StatsItem';
import type { CustomPropertyValue, TabData } from './types';

import classes from './Tabs.module.css';

interface TabSelectProps {
	data: TabData[];
}

interface TabPropertiesProps {
	label: string;
	choices: string[];
	data: CustomPropertyValue[];
}

interface LoadMoreButtonProps {
	label: string;
	searchParams: URLSearchParams;
}

const TabSelect = ({ data }: TabSelectProps) => {
	const [searchParams] = useSearchParams();

	return (
		<Tabs.Root className={classes.root} defaultValue={data[0]?.label}>
			<Tabs.List className={classes.list} aria-label="Select tab options">
				{data.map((tab) => (
					<Tabs.Trigger
						key={tab.label}
						className={classes.trigger}
						value={tab.label}
						aria-label={tab.label}
					>
						{tab.label}
					</Tabs.Trigger>
				))}
			</Tabs.List>
			{data.map((tab) => (
				<Tabs.Content key={tab.label} value={tab.label}>
					<div className={classes.items} data-empty={tab.items.length === 0}>
						{tab.items.map((item) => (
							<StatsItem key={item.label} tab={tab.label} {...item} />
						))}
						{tab.items.length === 0 && (
							<span className={classes.empty}>No records found...</span>
						)}
					</div>
					<LoadMoreButton label={tab.label} searchParams={searchParams} />
				</Tabs.Content>
			))}
		</Tabs.Root>
	);
};

const TabProperties = ({ label, choices, data }: TabPropertiesProps) => {
	const [searchParams] = useSearchParams();
	const { addFilter, getFilterEq, removeFilter } = useFilter();

	const choice = getFilterEq('prop_name') ?? '';
	const handleChoice = (value: string) => {
		// Replace any existing prop name filters if selected from combobox
		if (choice) {
			removeFilter('prop_name', 'eq', choice);
		}
		addFilter('prop_name', 'eq', value);
	};

	const items = data.map((item) => (
		<StatsItem
			key={item.name ?? item.value ?? 'Unknown'}
			tab={label}
			label={item.name ?? item.value ?? 'Unknown'}
			count={item.events}
		/>
	));

	return (
		<Tabs.Root className={classes.root} defaultValue={label}>
			<Tabs.List className={classes.list} aria-label="Select tab options">
				<Tabs.Trigger
					key={label}
					className={classes.trigger}
					value={label}
					aria-label={`Select ${label} tab`}
				>
					Properties
				</Tabs.Trigger>
			</Tabs.List>

			<Tabs.Content value={label}>
				<div>
					<Combobox
						root={{
							label: 'Select property',
							placeholder: 'Select property',
							emptyPlaceholder: 'No properties found...',
						}}
						search={{ placeholder: 'Search properties...' }}
						choices={choices}
						value={choice}
						setValue={handleChoice}
					/>
					<div className={classes.items} data-empty={data.length === 0}>
						{items}
						{data.length === 0 && (
							<span className={classes.empty}>No records found...</span>
						)}
					</div>
				</div>
				<LoadMoreButton label={label} searchParams={searchParams} />
			</Tabs.Content>
		</Tabs.Root>
	);
};

const LoadMoreButton = ({ label, searchParams }: LoadMoreButtonProps) => (
	<div className={classes['more-wrapper']}>
		<ButtonLink
			className={classes['more-button']}
			to={{
				pathname: `./${label.toLowerCase()}`,
				search: searchParams.toString(),
			}}
			prefetch="intent"
			aria-label={`Load more ${label} stats.`}
		>
			Load More
		</ButtonLink>
	</div>
);

export { TabSelect, TabProperties };
