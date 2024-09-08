import * as DropdownMenu from '@radix-ui/react-dropdown-menu';
import { ChevronDownIcon } from '@radix-ui/react-icons';
import { useMemo } from 'react';

import { ScrollArea } from '@/components/ScrollArea';

import classes from './WebsiteSelector.module.css';

interface WebsiteListComboboxProps {
	websites: string[];
	website: string;
	setWebsite: (website: string) => void;
}

export const WebsiteSelector = ({
	websites,
	website,
	setWebsite,
}: WebsiteListComboboxProps) => {
	const options = useMemo(
		() =>
			websites.map((value) => (
				<DropdownMenu.Item
					key={value}
					onSelect={() => setWebsite(value)}
					asChild
				>
					<button
						type="button"
						className={classes.item}
						aria-selected={value === website}
						data-active={value === website}
					>
						{value}
					</button>
				</DropdownMenu.Item>
			)),
		[websites, website, setWebsite],
	);

	return (
		<DropdownMenu.Root>
			<DropdownMenu.Trigger asChild>
				<button
					type="button"
					className={classes.trigger}
					aria-label="Select website hostname"
					disabled={websites.length === 0}
				>
					<span className={classes.label}>{website ?? 'No websites'}</span>
					<ChevronDownIcon />
				</button>
			</DropdownMenu.Trigger>

			<DropdownMenu.Content
				className={classes.dropdown}
				sideOffset={8}
				data-scroll={options.length > 5}
			>
				<ScrollArea vertical>{options}</ScrollArea>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	);
};
