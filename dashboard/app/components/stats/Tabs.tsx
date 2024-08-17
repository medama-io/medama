import * as Tabs from '@radix-ui/react-tabs';
import { useSearchParams } from '@remix-run/react';

import { ButtonLink } from '@/components/Button';

import { StatsItem } from './StatsItem';
import type { CustomEventValue, PageViewValue, TabData } from './types';

import classes from './Tabs.module.css';

interface TabSelectProps {
	data: TabData<PageViewValue>[];
}

interface TabPropertiesProps {
	data: TabData<CustomEventValue>[];
}

interface LoadMoreButtonProps {
	tab: TabData<PageViewValue | CustomEventValue>;
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
					<LoadMoreButton tab={tab} searchParams={searchParams} />
				</Tabs.Content>
			))}
		</Tabs.Root>
	);
};

const TabProperties = ({ data }: TabSelectProps) => {
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
					<LoadMoreButton tab={tab} searchParams={searchParams} />
				</Tabs.Content>
			))}
		</Tabs.Root>
	);
};

const LoadMoreButton = ({ tab, searchParams }: LoadMoreButtonProps) => (
	<div className={classes['more-wrapper']}>
		<ButtonLink
			className={classes['more-button']}
			to={{
				pathname: `./${tab.label.toLowerCase()}`,
				search: searchParams.toString(),
			}}
			prefetch="intent"
			aria-label={`Load more ${tab.label} stats.`}
		>
			Load More
		</ButtonLink>
	</div>
);

export { TabSelect, TabProperties };
