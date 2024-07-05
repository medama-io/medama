import { ActionIcon, Group, Tabs, Text, UnstyledButton } from '@mantine/core';
import { useDidUpdate, useMediaQuery } from '@mantine/hooks';
import { Link, useNavigate, useSearchParams } from '@remix-run/react';
import {
	DataTable,
	type DataTableColumn,
	type DataTableRowClickHandler,
	type DataTableSortStatus,
} from 'mantine-datatable';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { IconChevronLeft } from '@/components/icons/chevronleft';
import { IconChevronRight } from '@/components/icons/chevronright';
import { useFilter } from '@/hooks/use-filter';

import { formatCount, formatDuration, formatPercentage } from './formatter';
import type { DataRow, Dataset, Filter } from './types';
import { sortBy } from './utils';

import classes from './StatsTable.module.css';

type DataRowClick = DataTableRowClickHandler<DataRow>;

interface TablePaginationProps {
	page: number;
	pageSize: number;
	totalRecords: number;
	onPageChange: (page: number) => void;
	onPageSizeChange: (pageSize: (typeof PAGE_SIZES)[number]) => void;
}

interface StatsTableProps {
	query: Dataset;
	data: DataRow[];
}

interface QueryTableProps extends StatsTableProps {
	isMobile?: boolean;
}

const LABEL_MAP: Record<Dataset, string> = {
	pages: 'Pages',
	time: 'Time Spent',
	referrers: 'Referrers',
	sources: 'Sources',
	mediums: 'Mediums',
	campaigns: 'Campaigns',
	browsers: 'Browsers',
	os: 'Operating Systems',
	devices: 'Devices',
	countries: 'Countries',
	languages: 'Languages',
} as const;

const ACCESSOR_MAP: Record<Dataset, keyof DataRow> = {
	pages: 'path',
	time: 'duration',
	referrers: 'referrer',
	sources: 'source',
	mediums: 'medium',
	campaigns: 'campaign',
	browsers: 'browser',
	os: 'os',
	devices: 'device',
	countries: 'country',
	languages: 'language',
} as const;

const FILTER_MAP: Record<Dataset, Filter> = {
	pages: 'path',
	time: 'path',
	referrers: 'referrer',
	sources: 'utm_source',
	mediums: 'utm_medium',
	campaigns: 'utm_campaign',
	browsers: 'browser',
	os: 'os',
	devices: 'device',
	countries: 'country',
	languages: 'language',
} as const;

const PAGE_SIZES = [10, 25, 50, 100] as const;

// Preset columns
type PresetDataKeys =
	| 'path'
	| 'visitors'
	| 'visitors_percentage'
	| 'pageviews'
	| 'pageviews_percentage'
	| 'bounce_rate'
	| 'duration'
	| 'duration_lower_quartile'
	| 'duration_upper_quartile'
	| 'duration_percentage';

const PRESET_COLUMNS: Record<PresetDataKeys, DataTableColumn<DataRow>> = {
	path: { accessor: 'path', title: 'Path', width: '100%' },
	visitors: { accessor: 'visitors', title: 'Visitors', sortable: true },
	visitors_percentage: {
		accessor: 'visitors_percentage',
		title: 'Visitors %',
		render: (record) => formatPercentage(record.visitors_percentage ?? 0),
	},
	pageviews: {
		accessor: 'pageviews',
		title: 'Views',
		sortable: true,
		render: (record) => formatCount(record.pageviews ?? 0),
	},
	pageviews_percentage: {
		accessor: 'pageviews_percentage',
		title: 'Views %',
		render: (record) => formatPercentage(record.pageviews_percentage ?? 0),
	},
	bounce_rate: {
		accessor: 'bounce_rate',
		title: 'Bounce %',
		render: (record) =>
			formatPercentage((record.bounces ?? 0) / (record.visitors ?? 1)),
	},
	duration: {
		accessor: 'duration',
		title: 'Duration',
		sortable: true,
		render: (record) => formatDuration(record.duration ?? 0),
	},
	duration_lower_quartile: {
		accessor: 'duration_lower_quartile',
		title: 'Q1 (25%)',
		sortable: true,
		render: (record) => formatDuration(record.duration_lower_quartile ?? 0),
	},
	duration_upper_quartile: {
		accessor: 'duration_upper_quartile',
		title: 'Q3 (75%)',
		sortable: true,
		render: (record) => formatDuration(record.duration_upper_quartile ?? 0),
	},
	duration_percentage: {
		accessor: 'duration_percentage',
		title: 'Duration %',
		sortable: true,
		render: (record) => formatPercentage(record.duration_percentage ?? 0),
	},
};

const BackButton = () => {
	const [searchParams] = useSearchParams();

	return (
		<UnstyledButton
			component={Link}
			to={{
				pathname: '../',
				search: searchParams.toString(),
			}}
			className={classes.back}
		>
			<IconChevronLeft />
			<span>Go back</span>
		</UnstyledButton>
	);
};

const QueryTable = ({ query, data, isMobile }: QueryTableProps) => {
	// Pagination
	const [pageSize, setPageSize] = useState<(typeof PAGE_SIZES)[number]>(10);
	const [page, setPage] = useState(1);

	// Sorting
	const [sortStatus, setSortStatus] = useState<DataTableSortStatus<DataRow>>({
		columnAccessor: 'visitors',
		direction: 'desc',
	});

	const { isFilterActiveEq } = useFilter();
	const showLocales = isFilterActiveEq('language');
	const columns = useMemo(
		() => getColumnsForQuery(query, showLocales),
		[query, showLocales],
	);

	const records = useMemo(() => {
		const from = (page - 1) * pageSize;
		const to = from + pageSize;

		const sortedData = [...data].sort(
			sortBy(sortStatus.columnAccessor as keyof DataRow),
		);
		return sortStatus.direction === 'desc'
			? sortedData.reverse().slice(from, to)
			: sortedData.slice(from, to);
	}, [data, page, pageSize, sortStatus]);

	const handlePageChange = useCallback(
		(newPage: number) => {
			const maxPage = Math.ceil(data.length / pageSize);
			setPage(Math.max(1, Math.min(newPage, maxPage)));
		},
		[data.length, pageSize],
	);

	const handlePageSizeChange = useCallback(
		(newSize: (typeof PAGE_SIZES)[number]) => {
			setPageSize(newSize);
			setPage(1);
		},
		[],
	);

	const { addFilter } = useFilter();

	const handleFilter: DataRowClick = (row) => {
		const { record } = row;
		const filter = FILTER_MAP[query] ?? 'path';
		// Do not add filter if the page is on language and locales are already shown.
		if (filter !== 'language' || !showLocales) {
			const value =
				query === 'time'
					? record.path
					: record[ACCESSOR_MAP[query]] || 'Direct/None';
			addFilter(filter, 'eq', String(value));
		}
	};

	// Reset page when query or data length changes (from filters).
	useDidUpdate(() => {
		setPage(1);
		setPageSize(10);
	}, [query, data.length]);

	// Reset sort status when query changes.
	useDidUpdate(() => {
		setSortStatus({
			columnAccessor: 'visitors',
			direction: 'desc',
		});
	}, [query]);

	return (
		<div className={classes.tableWrapper}>
			<div className={classes.tableHeader}>
				<Text fz={14} fw={600} py={3}>
					{LABEL_MAP[query]}
				</Text>
			</div>
			{isMobile && <BackButton />}
			<DataTable
				classNames={{ header: classes.dataHeader }}
				minHeight={330}
				noRecordsText="No records found..."
				highlightOnHover
				withRowBorders={false}
				idAccessor={(record) => String(record[ACCESSOR_MAP[query]] ?? 'path')}
				records={records}
				columns={columns}
				sortStatus={sortStatus}
				onSortStatusChange={setSortStatus}
				onRowClick={handleFilter}
			/>
			<TablePagination
				page={page}
				pageSize={pageSize}
				totalRecords={data.length}
				onPageChange={handlePageChange}
				onPageSizeChange={handlePageSizeChange}
			/>
		</div>
	);
};

const TablePagination = ({
	page,
	pageSize,
	totalRecords,
	onPageChange,
	onPageSizeChange,
}: TablePaginationProps) => {
	const totalPages = Math.ceil(totalRecords / pageSize);

	return (
		<Group className={classes.pagination} px="lg" py="sm">
			<Group visibleFrom="sm">
				<span className={classes.viewspan}>View</span>
				{PAGE_SIZES.map((size) => (
					<ActionIcon
						key={size}
						variant="transparent"
						className={classes.pageSize}
						onClick={() => onPageSizeChange(size)}
						disabled={size === pageSize || totalRecords <= size}
						data-active={size === pageSize}
					>
						{size}
					</ActionIcon>
				))}
			</Group>
			<Group>
				<ActionIcon
					variant="transparent"
					className={classes.pageArrow}
					onClick={() => onPageChange(page - 1)}
					disabled={page <= 1}
				>
					<IconChevronLeft />
				</ActionIcon>
				<span>
					Page {page} of {totalPages}
				</span>
				<ActionIcon
					variant="transparent"
					className={classes.pageArrow}
					onClick={() => onPageChange(page + 1)}
					disabled={page >= totalPages}
				>
					<IconChevronRight />
				</ActionIcon>
			</Group>
		</Group>
	);
};

const getColumnsForQuery = (
	query: Dataset,
	filterActive?: boolean,
): DataTableColumn<DataRow>[] => {
	switch (query) {
		case 'pages':
			return [
				PRESET_COLUMNS.path,
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.pageviews,
				PRESET_COLUMNS.pageviews_percentage,
				PRESET_COLUMNS.bounce_rate,
				PRESET_COLUMNS.duration,
			];
		case 'time':
			return [
				PRESET_COLUMNS.path,
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.duration,
				PRESET_COLUMNS.duration_lower_quartile,
				PRESET_COLUMNS.duration_upper_quartile,
				PRESET_COLUMNS.duration_percentage,
			];
		case 'referrers':
		case 'sources':
		case 'mediums':
		case 'campaigns':
			return [
				{
					// Referrer, Source, Medium, Campaign
					accessor: ACCESSOR_MAP[query],
					title: LABEL_MAP[query].slice(0, -1),
					width: '100%',
					render: (record) => record[ACCESSOR_MAP[query]] || 'Direct/None',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_rate,
				PRESET_COLUMNS.duration,
			];
		case 'browsers':
		case 'devices':
			return [
				{
					// Browser, Device
					accessor: ACCESSOR_MAP[query],
					title: LABEL_MAP[query].slice(0, -1),
					width: '100%',
					render: (record) => record[ACCESSOR_MAP[query]] || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_rate,
				PRESET_COLUMNS.duration,
			];
		case 'languages':
			return [
				{
					// Browser, Device, Language
					accessor: ACCESSOR_MAP[query],
					title: filterActive ? 'Locale' : 'Language',
					width: '100%',
					render: (record) => record[ACCESSOR_MAP[query]] || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_rate,
				PRESET_COLUMNS.duration,
			];
		case 'os':
			return [
				{
					accessor: ACCESSOR_MAP[query],
					title: 'OS',
					width: '100%',
					render: (record) => record.os || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_rate,
				PRESET_COLUMNS.duration,
			];
		case 'countries':
			return [
				{
					accessor: ACCESSOR_MAP[query],
					title: 'Country',
					width: '100%',
					render: (record) => record.country || 'Unknown',
				},
				PRESET_COLUMNS.visitors,
				PRESET_COLUMNS.visitors_percentage,
				PRESET_COLUMNS.bounce_rate,
				PRESET_COLUMNS.duration,
			];
		default:
			throw new Error(`Invalid query: ${query}`);
	}
};

export const StatsTable = ({ query, data }: StatsTableProps) => {
	const navigate = useNavigate();
	const [searchParams] = useSearchParams();

	const handleTabChange = useCallback(
		(value: string | null) => {
			navigate(
				{
					pathname: `../${value}`,
					search: searchParams.toString(),
				},
				{ preventScrollReset: true },
			);
		},
		[navigate, searchParams],
	);

	// max-width: $mantine-breakpoint-lg
	const isMobile = useMediaQuery('(max-width: 75em)');

	return (
		<Tabs
			variant="unstyled"
			value={query}
			classNames={{
				root: classes.tabRoot,
				tab: classes.tab,
				list: classes.tabList,
				panel: classes.tabPanel,
			}}
			orientation="vertical"
			onChange={handleTabChange}
			keepMounted={false}
		>
			{!isMobile && (
				<Tabs.List>
					<div className={classes.listWrapper}>
						<BackButton />
						{Object.entries(LABEL_MAP).map(([key, label]) => (
							<Tabs.Tab key={key} value={key}>
								{label}
							</Tabs.Tab>
						))}
					</div>
				</Tabs.List>
			)}
			<Tabs.Panel value={query}>
				<QueryTable query={query} data={data} isMobile={isMobile} />
			</Tabs.Panel>
		</Tabs>
	);
};
