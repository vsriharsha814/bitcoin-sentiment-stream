import { Card } from "../../components/Card";
import { AlertTriangle, Bell } from "lucide-react";
import { Badge } from "../../components/Badge";

interface AlertsSectionProps {
  questions: string[];
  alerts: string[];
}

const AlertsSection = ({ questions, alerts }: AlertsSectionProps) => {
  return (
    <Card className="p-6 bg-black/20 border border-[#b829f7]/20 backdrop-blur-sm">
      <h2 className="text-2xl font-bold mb-4 text-transparent bg-clip-text bg-gradient-to-r from-[#b829f7] to-[#00fff9]">
        Questions & Alerts
      </h2>
      <div className="space-y-6">
        <div>
          <h3 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <Bell className="w-5 h-5 text-[#00fff9]" />
            Selected Questions
          </h3>
          <div className="flex flex-wrap gap-2">
            {questions.map((question) => (
              <Badge
                key={question}
                variant="secondary"
                className="bg-[#b829f7]/10 text-[#00fff9] hover:bg-[#b829f7]/20"
              >
                {question}
              </Badge>
            ))}
          </div>
        </div>
        <div>
          <h3 className="text-xl font-semibold mb-3 flex items-center gap-2">
            <AlertTriangle className="w-5 h-5 text-[#00fff9]" />
            Alert Conditions
          </h3>
          <div className="flex flex-wrap gap-2">
            {alerts.map((alert) => (
              <Badge
                key={alert}
                variant="secondary"
                className="bg-[#b829f7]/10 text-[#00fff9] hover:bg-[#b829f7]/20"
              >
                {alert}
              </Badge>
            ))}
          </div>
        </div>
      </div>
    </Card>
  );
};

export default AlertsSection;