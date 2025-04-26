import { Card } from "../../components/Card";
import { Bitcoin } from "lucide-react";

interface CryptoSectionProps {
  coins: string[];
}

const CryptoSection = ({ coins }: CryptoSectionProps) => {
  const getCoinIcon = (coin: string) => {
    switch (coin) {
      case "BTC":
        return <Bitcoin className="w-8 h-8 text-[#f7931a]" />;
      case "ETH":
        // Custom SVG for ETH since Eth isn't available in lucide-react
        return (
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            width="32" 
            height="32" 
            viewBox="0 0 256 417" 
            className="w-8 h-8 text-[#627eea]"
            fill="currentColor"
          >
            <path d="M127.9611 0.0369L125.1661 9.5894V285.168L127.9611 288.342L255.9231 212.175L127.9611 0.0369Z" />
            <path d="M127.962 0.0371L0 212.175L127.962 288.343V154.49V0.0371Z" />
            <path d="M127.9609 312.1866L126.3859 313.9956V409.8326L127.9609 414.8806L255.9999 236.5786L127.9609 312.1866Z" />
            <path d="M127.962 414.8808V312.1858L0 236.5788L127.962 414.8808Z" />
            <path d="M127.9609 288.3431L255.9229 212.1761L127.9609 154.4921V288.3431Z" />
            <path d="M0.0001 212.176L127.9611 288.343V154.492L0.0001 212.176Z" />
          </svg>
        );
      default:
        return null;
    }
  };

  return (
    <Card className="p-6 bg-black/20 border border-[#b829f7]/20 backdrop-blur-sm">
      <h2 className="text-2xl font-bold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-[#b829f7] to-[#00fff9]">
        Tracked Coins
      </h2>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {coins.map((coin) => (
          <div
            key={coin}
            className="flex items-center space-x-3 p-4 rounded-lg bg-black/40 border border-[#b829f7]/10 hover:border-[#b829f7]/30 transition-all"
          >
            {getCoinIcon(coin)}
            <span className="font-bold">{coin}</span>
          </div>
        ))}
      </div>
    </Card>
  );
};

export default CryptoSection;