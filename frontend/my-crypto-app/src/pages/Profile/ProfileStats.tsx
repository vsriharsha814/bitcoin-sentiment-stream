import { Card } from "../../components/Card";

interface ProfileStatsProps {
  user: {
    coins: string[];
    questions: string[];
    alerts: string[];
  };
}

const ProfileStats = ({ user }: ProfileStatsProps) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Card className="p-6 text-center bg-black/20 border border-[#b829f7]/20 backdrop-blur-sm hover:border-[#b829f7]/40 transition-all">
        <h3 className="text-2xl font-bold text-[#00fff9]">{user.coins.length}</h3>
        <p className="text-gray-400">Coins Tracked</p>
      </Card>
      <Card className="p-6 text-center bg-black/20 border border-[#b829f7]/20 backdrop-blur-sm hover:border-[#b829f7]/40 transition-all">
        <h3 className="text-2xl font-bold text-[#00fff9]">{user.questions.length}</h3>
        <p className="text-gray-400">Active Questions</p>
      </Card>
      <Card className="p-6 text-center bg-black/20 border border-[#b829f7]/20 backdrop-blur-sm hover:border-[#b829f7]/40 transition-all">
        <h3 className="text-2xl font-bold text-[#00fff9]">{user.alerts.length}</h3>
        <p className="text-gray-400">Alert Conditions</p>
      </Card>
    </div>
  );
};

export default ProfileStats;